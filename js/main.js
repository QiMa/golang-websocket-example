import React from 'react';
import { render } from 'react-dom';

var ws = new WebSocket("ws://localhost:4000/ws");
ws.onopen = function () {
    console.log('Connection is open');
};

let Chat = React.createClass({
    getInitialState: function () {
        return {
            messages: [],
        }
    },
    componentDidMount: function () {
        ws.onmessage = this.handleWebSocketMessage;
    },
    componentWillUnmount: function () {
        ws.onmessage = function () {};
    },
    handleWebSocketMessage: function (event) {
        this.setState({
            messages: [event.data].concat(this.state.messages)
        });
    },
    handleKeyUp: function (event) {
        if (event.keyCode === 13) {
            var message = event.target.value;
            event.target.value = '';
            ws.send(message);
        }
    },
    render: function () {
        return <div>
            <p>
                <input type="text" onKeyUp={this.handleKeyUp} />
            </p>
            {this.state.messages.map(function (message, index) {
                return <p key={index} >{message}</p>;
            })}
        </div>;
    }
});

render(<Chat/>, document.getElementById('app'));
