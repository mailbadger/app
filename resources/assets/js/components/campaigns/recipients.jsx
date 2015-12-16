
import React, {Component} from 'react';
import List from '../../entities/list.js';

const l = new List();

export default class Recipients extends Component {

    constructor(props) {
        super(props);
        this.state = {
            lists: [],
            total: 0
        };
    }

    componentDidMount() {
        l.all().done((res) => {
            this.setState({lists: res});
        });
    }

    render() {
        let listOpts = (list) => {
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
}
