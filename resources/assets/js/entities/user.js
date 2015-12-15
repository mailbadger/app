
export default class User {
   
    getSettings() {
        return $.ajax({
            type: 'GET',
            url: url_base + '/dashboard/user-settings'
        });
    }

    saveSettings(data) {
        return $.ajax({
            type: 'POST',
            url: url_base + '/dashboard/settings',
            data: data,
            dataType: 'json'
        });
    }
}
