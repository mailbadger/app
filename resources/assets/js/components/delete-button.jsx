
import sweetalert from 'sweetalert';
import React, {Component} from 'react';

export default class DeleteButton extends Component {

    constructor(props) {
        super(props);

        this.handleSubmit = this.handleSubmit.bind(this);
        this.deleteEntity = this.deleteEntity.bind(this);
    }

    deleteEntity() { 
        this.props.entity.delete(this.props.resourceId)
        .done(() => {
            sweetalert({
                title: "Success",
                text: "The entity was successfully removed!",
                type: "success"
            }, () => {
                this.props.success();
            });
        })
        .fail(() => {
            sweetalert('Could not delete', 'Could not delete the entity. Try again.', 'error');
        });
    }

    handleSubmit(e) {
        e.preventDefault();
        sweetalert({
                title: "Are you sure?",
                text: "You will not be able to recover this!",
                type: "warning",
                showCancelButton: true,
                confirmButtonColor: "#DD6B55",
                confirmButtonText: "Yes, delete it!",
                closeOnConfirm: false
            }, this.deleteEntity);
    }

    render() {
        return (
            <form onSubmit={this.handleSubmit}>
                <input type="hidden" name="_method" value="DELETE"/>
                <button type="submit"><span className="glyphicon glyphicon-trash"></span></button>
            </form>
        );
    }
}
