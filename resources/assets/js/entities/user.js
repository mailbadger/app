
var User = function () {
   
    this.getSettings = function () {
        return $.ajax({
            type: 'GET',
            url: url_base + '/dashboard/user-settings'
        });
    };

    this.saveSettings = function (data) {
        return $.ajax({
            type: 'POST',
            url: url_base + '/dashboard/settings',
            data: data,
            dataType: 'json'
        });
    };
}

module.exports = User;
