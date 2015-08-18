/**
 * Created by filip on 18.8.15.
 */
var Field = function () {

    var url = url_base + '/api/fields';

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

};

module.exports = Field;