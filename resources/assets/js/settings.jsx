
import React from 'react';
import SettingsForm from './components/settings-form.jsx';
import User from './entities/user.js';

const u = new User();

u.getSettings().done((res) => {  
    React.render(<SettingsForm data={res} />, document.getElementById('settings'));
});

