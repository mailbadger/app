
import * as bootpag from 'bootpag/lib/jquery.bootpag.min.js';

import React, {Component} from 'react';
import DeleteButton from '../delete-button.jsx';
import Subscriber from '../../entities/subscriber.js';

class SubscriberRow extends Component {
    render() {
        return (
            <tr>
                <td>{this.props.data.name}</td>
                <td>{this.props.data.email}</td>
                <td>
                    <DeleteButton success={this.props.handleDelete} entity={this.props.entity} resourceId={this.props.data.id} />
                </td>
            </tr>
        );
    }
}

export default class SubscribersTable extends Component {

    constructor(props) {
        super(props);

        this.s = new Subscriber(this.props.listId);

        this.state = {
            subscribers: {
                data: []
            }
        };

        this.handleDelete = this.handleDelete.bind(this);
    }
    
    getAllSubscribers(component) {
        let data = {
            paginate: true,
            per_page: 10,
            page: 1
        };

        this.s.all(data).done((res) => {
            this.setState({subscribers: res});

            $('.pagination').bootpag({
                total: res.last_page,
                page: res.current_page,
                maxVisible: 5
            }).on("page", (event, num) => {
                data.page = num;
                s.all(data).done((res) => {
                    this.setState({subscribers: res});
                    $('.pagination').bootpag({page: res.current_page});
                });
            });
        });
    }

    componentDidMount() {
        getAllSubscribers();
    }

    handleDelete() {
        getAllSubscribers();
    }

    render() {
        let rows = (data) => {
            return <SubscriberRow key={data.id} entity={this.s} handleDelete={this.handleDelete} listId={this.props.listId} data={data}/>
        };

        return (
            <div>
                <table className="table table-responsive table-striped table-hover">
                    <thead>
                    <tr>
                        <th>Subscriber name</th>
                        <th>Email</th>
                        <th>Delete</th>
                    </tr>
                    </thead>
                    <tbody>
                    {this.state.subscribers.data.map(rows)}
                    </tbody>
                </table>
                <div className="col-lg-12 pagination text-center"></div>
            </div>
        );
    }
}
