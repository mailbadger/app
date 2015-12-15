/**
 * Created by filip on 22.7.15.
 */
import Entity from './entity'

export default class Campaign extends Entity {
    
    constructor() {
        super('campaigns');
    }

    testSend(emails, id) {
        return $.ajax({
            type: 'POST',
            url: this.url + '/' + id + '/test-send',
            data: {
                emails: emails,
                id: id
            },
            dataType: 'json'
        });
    }

    send(lists, id) {
        return $.ajax({
            type: 'POST',
            url: this.url + '/' + id + '/send',
            data: {
                lists: lists,
                id: id
            },
            dataType: 'json'
        });
    }
}
