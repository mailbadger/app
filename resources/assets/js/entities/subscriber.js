
import Entity from './entity';

export default class Subscriber extends Entity {

    constructor(listId) {
        super('lists/' + listId + '/subscribers');
    }
}
