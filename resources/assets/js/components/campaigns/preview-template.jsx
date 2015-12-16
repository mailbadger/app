
import React, {Component} from 'react';
import Template from '../../entities/template.js';

const t = new Template();

export default class PreviewTemplate extends Component {
    render() {
        return (
            <div>
                <h3>Preview template</h3>
                <blockquote>
                    <p><strong>From</strong> <span className="label label-default">{this.props.from}</span></p>
                </blockquote>

                <blockquote>
                    <p><strong>Subject</strong> <span className="label label-default">{this.props.subject}</span></p>
                </blockquote>

                <iframe id="preview" src={url_base + '/api/templates/content/' + this.props.tid} className="col-lg-12" />
            </div>
        )
    }
}
