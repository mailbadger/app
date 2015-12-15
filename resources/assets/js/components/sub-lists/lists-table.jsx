
import * as bootpag from 'bootpag/lib/jquery.bootpag.min.js';

import React, {Component} from 'react';
import DeleteButton from '../delete-button.jsx';
import List from '../../entities/list.js';

const l = new List();

const getAllLists = (component) => {
    let data = {
        paginate: true,
        per_page: 10,
        page: 1
    };

    l.all(data).done((res) => {
        component.setState({lists: res});
        $('.pagination').bootpag({
            total: res.last_page,
            page: res.current_page,
            maxVisible: 5
        }).on("page", (event, num) => {
            data.page = num;
            l.all(data).done((res) =>{
                component.setState({lists: res});
                $('.pagination').bootpag({page: res.current_page});
            });
        });
    });
}

class ListRow extends Component {

    constructor(props) {
        super(props);

        this.showList = this.showList.bind(this);
        this.editList = this.editList.bind(this);
    }

    showList() {
        this.props.showList(this.props.data.id);
    }

    editList() {
        this.props.editList(this.props.data.id);
    }

    render() {
        return (
            <tr>
                <td><a href="#" onClick={this.showList}>{this.props.data.name}</a></td>
                <td>{this.props.data.total_subscribers}</td>
                <td>
                    <a href="#" onClick={this.editList}><span className="glyphicon glyphicon-pencil"></span></a>
                </td>
                <td>
                    <DeleteButton success={this.props.handleDelete} entity={l} resourceId={this.props.data.id} />
                </td>
            </tr>
        );
    }
}

export default class ListsTable extends Component {

    constructor(props) {
        super(props);
        this.state = {
            lists: {
                data: []
            }
        };

        this.handleDelete = this.handleDelete.bind(this);
    }

    componentDidMount() {
        getAllLists(this);
    }

    handleDelete() {
        getAllLists(this);
    }

    render() {
        let rows = (data) => {
            return <ListRow key={data.id} data={data} handleDelete={this.handleDelete} showList={this.props.showList}
                            editList={this.props.editList}/>
        };

        return (
            <div>
                <table className="table table-responsive table-striped table-hover">
                    <thead>
                    <tr>
                        <th>List name</th>
                        <th>Subscribers</th>
                        <th>Edit</th>
                        <th>Delete</th>
                    </tr>
                    </thead>
                    <tbody>
                    {this.state.lists.data.map(rows)}
                    </tbody>
                </table>
                <div className="col-lg-12 pagination text-center"></div>
            </div>
        );
    }
}
