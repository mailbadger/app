
import Entity from './entity';

export default class Field extends Entity {

    constructor(listId) {
        super('lists/' + listId + '/fields');
    }
}
