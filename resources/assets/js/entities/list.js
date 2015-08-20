/**
 * Created by filip on 3.8.15.
 */
/**
 * Created by filip on 25.7.15.
 */
var List = function () {

    var url = url_base + '/api/lists';

    this.all = function (paginate, perPage, page) {
        return $.ajax({
            type: 'GET',
            url: url,
            data: {
                paginate: paginate,
                per_page: perPage,
                page: page
            },
            dataType: 'json'
        });
    };

    this.get = function (id) {
        return $.ajax({
            type: 'GET',
            url: url + '/' + id,
            dataType: 'json'
        });
    };

    this.delete = function (id) {
        return $.ajax({
            type: 'DELETE',
            url: url + '/' + id,
            dataType: 'json'
        });
    };

    this.create = function (data) {
        return $.ajax({
            type: 'POST',
            url: url,
            data: data,
            dataType: 'json'
        });
    };

    this.update = function (data, id) {
        return $.ajax({
            type: 'PUT',
            url: url + '/' + id,
            data: data,
            dataType: 'json'
        });
    };

    this.allSubscribers = function (listId, paginate, perPage, page) {
        return $.ajax({
            type: 'GET',
            url: url + '/' + listId + '/subscribers',
            data: {
                paginate: paginate,
                per_page: perPage,
                page: page
            },
            dataType: 'json'
        });
    };

    this.createSubscribers = function (listId, file) {
        var data = new FormData();
        data.append('subscribers', file);
        data.append('list_id', listId);

        return $.ajax({
            type: 'POST',
            url: url + '/' + listId + '/import-subscribers',
            data: data,
            processData: false,
            contentType: false,
            cache: false,
            dataType: 'json'
        });
    };

    this.deleteSubscribers = function (listId, file) {
        //TODO add delete subscribers functionality
    };

    this.deleteSubscriber = function (listId, id) {
        return $.ajax({
            type: 'DELETE',
            url: url + '/' + listId + '/subscribers/' + id,
            dataType: 'json'
        });
    };

    this.allFields = function (listId, paginate, perPage, page) {
        return $.ajax({
            type: 'GET',
            url: url + '/' + listId + '/fields',
            data: {
                paginate: paginate,
                per_page: perPage,
                page: page
            },
            dataType: 'json'
        });
    };

    this.createField = function (listId, data) {
        return $.ajax({
            type: 'POST',
            url: url + '/' + listId + '/fields',
            data: data,
            dataType: 'json'
        });
    };

    this.updateField = function (listId, data, id) {
        return $.ajax({
            type: 'PUT',
            url: url + '/' + listId + '/fields/' + id,
            data: data,
            dataType: 'json'
        });
    };

    this.deleteField = function (listId, id) {
        return $.ajax({
            type: 'DELETE',
            url: url + '/' + listId + '/fields/' + id,
            dataType: 'json'
        });
    };
};

module.exports = List;