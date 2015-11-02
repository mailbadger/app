/**
 * Created by filip on 22.7.15.
 */
var Campaign = function () {

    var url = url_base + '/api/campaigns';

    this.all = function (paginate, perPage, page, filters) {
        return $.ajax({
            type: 'GET',
            url: url,
            data: {
                paginate: paginate,
                per_page: perPage,
                page: page,
                filters: filters
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

    this.testSend = function (emails, id) {
        return $.ajax({
            type: 'POST',
            url: url + '/' + id + '/test-send',
            data: {
                emails: emails,
                id: id
            },
            dataType: 'json'
        });
    };

    this.send = function(lists, id) {
        return $.ajax({
            type: 'POST',
            url: url + '/' + id + '/send',
            data: {
                lists: lists,
                id: id
            },
            dataType: 'json'
        });
    }
};

module.exports = Campaign;
