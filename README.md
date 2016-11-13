[![Build Status](https://travis-ci.org/FilipNikolovski/news-maily.svg?branch=golang)](https://travis-ci.org/FilipNikolovski/news-maily)
[![GitHub license](https://img.shields.io/badge/license-Apache%202-blue.svg)](https://raw.githubusercontent.com/FilipNikolovski/news-maily/master/LICENSE.md)

# news-maily

Self hosted newsletter mail system written in go.

# TODO

- [x] Add login and logout actions.
- [x] Refactor jwt secret key func.
- [x] Add template actions - get template(s), add template, delete template.
- [x] Add campaign storage functions - add, get, delete etc.
- [x] Add campaign actions.
- [x] Add lists storage functions.
- [x] Add lists actions.
- [x] Add subscribers to list.
- [x] Delete subscribers from list.
- [ ] Unsubscribe functionality.
- [ ] Test send email functionality.
- [ ] Send campaign (should be done as a background process).
- [ ] Schedule campaigns.
- [ ] Track emails.
- [ ] Track link clicks.
- [ ] Change password.
- [ ] Forgot password.
- [ ] Account settings.
- [ ] Bounces/Complaints webhook (for now it would only work with SES).
- [ ] Separate subscribers whose emails have bounced or those that have complained (non-existent email addresses, or marked as spam).
- [ ] Do not include 'deleted', blacklisted or opted-out subscribers in the sending campaign process.
- [ ] Export subscribers.
- [ ] Import subscribers from csv/excel feature.
- [ ] Mass delete subscribers via csv/excel file (soft delete).
- [ ] Environment vars.
- [ ] Create docker image for easier installation.
- [ ] Proper README, website and logo.
- [ ] Action tests.
