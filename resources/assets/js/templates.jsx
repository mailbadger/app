
import React, {Component} from 'react';
import TemplateForm from './components/templates/template-form.jsx';
import TemplatesTable from './components/templates/templates-table.jsx';
import CreateNewButton from './components/create-new-button.jsx';
import Template from './entities/template.js';

const t = new Template();

class Templates extends Component {

    constructor(props) {
        super(props);
        this.state = {
            step: '',
            campaign: {}
        };

        this.editTemplate = this.editTemplate.bind(this);
        this.back = this.back.bind(this);
    }

    editTemplate(id) {
        t.get(id).done((res) => {
            this.setState({step: 'edit', template: res});
        });
    }

    back() {
        this.setState({step: ''});
    }
    
    render() {
        switch (this.state.step) {
            case 'edit':
                return <TemplateForm data={this.state.template} edit={true} back={this.back}/>;
            default:
                return (
                    <div>
                        <CreateNewButton url={url_base + '/dashboard/new-template'} text="Create new template"/>

                        <div className="row">
                            <TemplatesTable editTemplate={this.editTemplate}/>
                        </div>
                    </div>
                );
        }
    }
}

React.render(<Templates />, document.getElementById('templates'));
