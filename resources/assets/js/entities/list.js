/**
 * Created by filip on 3.8.15.
 */

import Entity from './entity';

export default class List extends Entity {

    constructor() {
        super('lists');
    }

    createSubscribers(listId, file) {
        const data = new FormData();
        data.append('subscribers', file);
        data.append('list_id', listId);

        return $.ajax({
            type: 'POST',
            url: this.url + '/' + listId + '/import-subscribers',
            data: data,
            processData: false,
            contentType: false,
            cache: false,
            dataType: 'json'
        });
    }

    deleteSubscribers(listId, file) {
        const data = new FormData();
        data.append('subscribers', file);
        data.append('list_id', listId);

        return $.ajax({
            type: 'POST',
            url: this.url + '/' + listId + '/mass-delete-subscribers',
            data: data,
            processData: false,
            contentType: false,
            cache: false,
            dataType: 'json'
        });
    }
}
