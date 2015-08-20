/** @jsx React.DOM */

var React = require('react');
var SubscribersTable = require('./subscribers-table.jsx');
var AddSubscribers = require('./add-subscribers.jsx');
var List = require('../../entities/list.js');

var l = new List();

var SubscriberButtons = React.createClass({
    render: function () {
        return (
            <div className="row">
                <div className="col-lg-6">
                    <div className="row">
                        <button className="btn col-lg-3" onClick={this.props.addSubscribers}><span
                            className="glyphicon glyphicon-save-file"></span> Add
                            subscribers
                        </button>
                        <button className="btn col-lg-3 col-lg-offset-1"><span
                            className="glyphicon glyphicon-trash"></span> Delete
                            subscribers
                        </button>
                        <button className="btn col-lg-3 col-lg-offset-1"><span
                            className="glyphicon glyphicon-export"></span> Export
                            subscribers
                        </button>
                    </div>
                </div>
            </div>
        )
    }
});

var ListInfo = React.createClass({
    editList: function () {
        this.props.editList(this.props.list.id);
    },
    customFields: function () {
        this.props.customFields(this.props.list.id);
    },
    render: function () {
        return (
            <div className="row" style={{marginTop: '20px'}}>
                <div className="col-lg-12 well">
                    <div className="pull-left">
                        List: <span className="label label-primary">{this.props.list.name}</span><span> | </span>
                        <a href="#" onClick={this.props.back}>Back to lists</a>
                    </div>
                    <div className="col-lg-6 pull-right">
                        <a href="#" onClick={this.editList} className="pull-right col-lg-offset-1"><span
                            className="glyphicon glyphicon-wrench"></span> Edit list</a>
                        <a href="#" onClick={this.customFields} className="pull-right"><span
                            className="glyphicon glyphicon-list-alt"></span> Custom
                            fields</a>
                    </div>
                </div>
            </div>
        )
    }
});

var ListView = React.createClass({
    render: function () {
        return (
            <div>
                <SubscriberButtons addSubscribers={this.props.addSubscribers}/>
                <ListInfo list={this.props.list} editList={this.props.editList} customFields={this.props.customFields}
                          back={this.props.back}/>
                <SubscribersTable listId={this.props.list.id}/>
            </div>
        );
    }
});

var SubscribersList = React.createClass({
    addSubscribers: function () {
        this.setState({component: <AddSubscribers listId={this.props.list.id} back={this.back}/>});
    },
    back: function () {
        this.setState({
            component: <ListView addSubscribers={this.addSubscribers} list={this.props.list}
                                 editList={this.props.editList} customFields={this.props.customFields}
                                 back={this.props.back}/>
        });
    },
    getInitialState: function () {
        return {
            component: <ListView addSubscribers={this.addSubscribers} list={this.props.list}
                                 editList={this.props.editList} customFields={this.props.customFields}
                                 back={this.props.back}/>
        };
    },
    render: function () {
        return this.state.component;
    }
});

module.exports = SubscribersList;