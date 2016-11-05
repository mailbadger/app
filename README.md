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
- [ ] Add lists storage functions.
- [ ] Add lists actions.
- [ ] Import subscribers from csv/excel feature.
- [ ] Create and edit list fields.
- [ ] Subscriber import file proper validation by fields = columns.
- [ ] Mass delete subscribers via csv/excel file.
- [ ] Test send email functionality.
- [ ] Send campaign (should be done as a background process).
- [ ] Schedule campaigns.
- [ ] Export subscribers.
- [ ] Track emails route.
- [ ] Account settings.
- [ ] Bounces/Complaints webhook (for now it would only work with SES).
- [ ] Environment vars.
- [ ] Create docker image for easier installation.
- [ ] Proper README, website and logo.
- [ ] Campaign http tests
- [ ] Template http tests
- [ ] User http tests
