/** @jsx React.DOM */

require('sweetalert');
var React = require('react');
var List = require('../../entities/list.js');
var l = new List();

var Recipients = React.createClass({
    getInitialState: function () {
        return {lists: [], total: 0};
    },
    componentDidMount: function () {
        l.all().done(function (res) {
            this.setState({lists: res});
        }.bind(this));
    },
    render: function () {
        var listOpts = function (list) {
            return <option key={list.id} value={list.id}>{list.name}</option>
        };
        return (
            <div>
                <h3>Select subscribers</h3>

                <div className="form-group">
                    <label htmlFor="subscribers">Select subscriber list(s)</label>
                    <select className="form-control" ref="subscribers" id="subscribers" multiple={true} onChange={this.handleChange} required>
                        {this.state.lists.map(listOpts)}
                    </select> 
                </div>
            </div>
        )
    }
});

module.exports = Recipients;
