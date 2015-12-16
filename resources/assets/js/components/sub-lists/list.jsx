
import React, {Component} from 'react';
import SubscribersTable from './subscribers-table.jsx';
import AddSubscribers from './add-subscribers.jsx';
import DeleteSubscribers from './delete-subscribers.jsx';
import List from '../../entities/list.js';

const l = new List();

class SubscriberButtons extends Component {
    render() {
        return (
            <div className="row">
                <div className="col-lg-6">
                    <div className="row">
                        <a className="btn btn-default col-lg-3" onClick={this.props.addSubscribers}><span
                            className="glyphicon glyphicon-save-file"></span> Add
                            subscribers
                        </a>
                        <a className="btn btn-default col-lg-3 col-lg-offset-1"
                           onClick={this.props.deleteSubscribers}>
                            <span className="glyphicon glyphicon-trash"></span> Delete subscribers
                        </a>
                        <a href={url_base + '/api/lists/' + this.props.listId + '/export-subscribers'}
                           className="btn btn-default col-lg-3 col-lg-offset-1"><span
                            className="glyphicon glyphicon-export"></span> Export
                            subscribers
                        </a>
                    </div>
                </div>
            </div>
        )
    }
}

class ListInfo extends Component {

    constructor(props) {
        super(props);

        this.editList = this.editList.bind(this);
        this.customFields = this.customFields.bind(this);
    }

    editList() {
        this.props.editList(this.props.list.id);
    }

    customFields() {
        this.props.customFields(this.props.list.id);
    }

    render() {
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
}

class ListView extends Component {
    render() {
        return (
            <div>
                <SubscriberButtons listId={this.props.list.id} addSubscribers={this.props.addSubscribers}
                                   deleteSubscribers={this.props.deleteSubscribers}/>
                <ListInfo list={this.props.list} editList={this.props.editList} customFields={this.props.customFields}
                          back={this.props.back}/>
                <SubscribersTable listId={this.props.list.id}/>
            </div>
        );
    }
}

export default class SubscribersList extends Component {

    constructor(props) {
        super(props);

        this.addSubscribers = this.addSubscribers.bind(this);
        this.deleteSubscribers = this.deleteSubscribers.bind(this);
    }

    addSubscribers() {
        this.setState({component: <AddSubscribers listId={this.props.list.id} back={this.back}/>});
    }

    deleteSubscribers() {
        this.setState({component: <DeleteSubscribers listId={this.props.list.id} back={this.back}/>});
    }

    back() {
        this.setState({
            component: <ListView addSubscribers={this.addSubscribers} deleteSubscribers={this.deleteSubscribers}
                                 list={this.props.list} editList={this.props.editList}
                                 customFields={this.props.customFields} back={this.props.back}/>
        });
    }

    getInitialState() {
        return {
            component: <ListView addSubscribers={this.addSubscribers} deleteSubscribers={this.deleteSubscribers}
                                 list={this.props.list} editList={this.props.editList}
                                 customFields={this.props.customFields} back={this.props.back}/>
        };
    }

    render() {
        return this.state.component;
    }
}
