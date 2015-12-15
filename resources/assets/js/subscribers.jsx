
import React, {Component} from 'react';

import SubscribersList from './components/sub-lists/list.jsx';
import ListForm from './components/sub-lists/list-form.jsx';
import CustomFields from './components/sub-lists/custom-fields.jsx';
import ListsTable from './components/sub-lists/lists-table.jsx';
import CreateNewButton from './components/create-new-button.jsx';
import List from './entities/list.js';

const l = new List();

class Lists extends Component {

    constructor(props) {
        super(props);
        this.state = {
            step: '',
            list: {}
        };

        this.showList = this.showList.bind(this);
        this.editList = this.editList.bind(this);
        this.customFields = this.customFields.bind(this);
        this.back = this.back.bind(this);
    }
 
    showList(id) {
        l.get(id).done((res) => {
            this.setState({step: 'show', list: res});
        });
    }

    editList(id) {
        l.get(id).done((res) => {
            this.setState({step: 'edit', list: res});
        });
    }
    
    customFields(id) {
        this.setState({step: 'custom-fields'});
    }

    back() {
        this.setState({step: ''});
    }

    render() {
        switch (this.state.step) {
            case 'show':
                return <SubscribersList list={this.state.list} editList={this.editList} customFields={this.customFields}
                    back={this.back}/>;
            case 'edit':
                return <ListForm data={this.state.list} edit={true} back={this.back}/>;
            case 'custom-fields':
                return <CustomFields listId={this.state.list.id} back={this.back}/>;
            default:
                return (
                    <div>
                        <CreateNewButton url={url_base + '/dashboard/new-subscribers'} text="Create new list"/>

                        <div className="row">
                            <ListsTable showList={this.showList} editList={this.editList}/>
                        </div>
                    </div>
                );     
        }
    }
}

React.render(<Lists />, document.getElementById('sub-lists'));
