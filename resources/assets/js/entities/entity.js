
export default class Entity {
    
    constructor(resource) {
        this.url = url_base + '/api/' + resource;
    }

    all(data) { 
        return $.ajax({
            type: 'GET',
            url: this.url,
            data: data,
            dataType: 'json'
        });
    }

    get(id) {
        return $.ajax({
            type: 'GET',
            url: this.url + '/' + id,
            dataType: 'json'
        });
    }

    delete(id) {
        return $.ajax({
            type: 'DELETE',
            url: this.url + '/' + id,
            dataType: 'json'
        });
    };

    create(data) {
        return $.ajax({
            type: 'POST',
            url: this.url,
            data: data,
            dataType: 'json'
        });
    };

    update(data, id) {
        return $.ajax({
            type: 'PUT',
            url: this.url + '/' + id,
            data: data,
            dataType: 'json'
        });
    } 
}
