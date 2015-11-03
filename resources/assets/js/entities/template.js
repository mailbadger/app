/**
 * Created by filip on 25.7.15.
 */
var Template = function () {

    var url = url_base + '/api/templates';

    this.all = function (data) {
        return $.ajax({
            type: 'GET',
            url: url,
            data: data,
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

    this.update = function(data, id) {
        return $.ajax({
            type: 'PUT',
            url: url + '/' + id,
            data: data,
            dataType: 'json'
        });
    }
};

module.exports = Template;
