/** @jsx React.DOM */

var React = require('react');

var ListForm = require('./components/sub-lists/list-form.jsx');
var ListsTable = require('./components/sub-lists/lists-table.jsx');
var CreateNewButton = require('./components/create-new-button.jsx');
var List = require('./entities/list.js');

var l = new List();

var Lists = React.createClass({
    getInitialState: function () {
        return {step: '', template: {}}
    },
    editList: function (id) {
        l.get(id).done(function (res) {
            this.setState({step: 'edit', template: res});
        }.bind(this));
    },
    back: function () {
        this.setState({step: ''});
    },
    render: function () {
        switch (this.state.step) {
            case 'edit':
                return <ListForm data={this.state.template} edit={true} back={this.back}/>;
            default:
                return (
                    <div>
                        <CreateNewButton url={url_base + '/dashboard/new-subscribers'} text="Create new list"/>

                        <div className="row">
                            <ListsTable editList={this.editList}/>
                        </div>
                    </div>
                );
        }
    }
});

React.render(<Lists />, document.getElementById('sub-lists'));