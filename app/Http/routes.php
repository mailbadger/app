<?php

/*
|--------------------------------------------------------------------------
| Application Routes
|--------------------------------------------------------------------------
 */

//No csrf middleware on these routes
Route::post('api/emails/bounces', 'Api\EmailController@bounces');

Route::post('api/emails/complaints', 'Api\EmailController@complaints');

Route::get('api/emails/opens', 'Api\EmailController@opens');

//Api routes
Route::group(['middleware' => 'csrf'], function () {
    // Authentication routes
    Route::get('auth/login', 'Auth\AuthController@getLogin');
    Route::post('auth/login', 'Auth\AuthController@postLogin');
    Route::get('auth/logout', 'Auth\AuthController@getLogout');

    // Password reset link request routes
    Route::get('password/email', 'Auth\PasswordController@getEmail');
    Route::post('password/email', 'Auth\PasswordController@postEmail');

    // Password reset routes
    Route::get('password/reset/{token}', 'Auth\PasswordController@getReset');
    Route::post('password/reset', 'Auth\PasswordController@postReset');

    //Dashboard controller
    Route::controller('dashboard', 'DashboardController');

    //Api controllers
    Route::group(['namespace' => 'Api', 'prefix' => 'api'], function () {
        Route::post('campaigns/{id}/send', 'CampaignController@send');

        Route::post('campaigns/{id}/test-send', 'CampaignController@testSend');

        Route::resource('campaigns', 'CampaignController', ['except' => ['create', 'edit']]);

        Route::get('templates/content/{id}', 'TemplateController@showContent');

        Route::resource('templates', 'TemplateController', ['except' => ['create', 'edit']]);

        Route::post('lists/{id}/import-subscribers', 'ListsSubscriberController@import');

        Route::post('lists/{id}/mass-delete-subscribers', 'ListsSubscriberController@massDelete');

        Route::get('lists/{id}/export-subscribers', 'ListsSubscriberController@export');

        Route::resource('lists', 'ListsController', ['except' => ['create', 'edit']]);

        Route::resource('lists.subscribers', 'ListsSubscriberController', ['except' => ['create', 'edit']]);

        Route::resource('lists.fields', 'ListsFieldController', ['except' => ['create', 'edit']]);

    });
});

Route::get('/', [
    'middleware' => 'guest',
    function () {
        return view('auth.login');
    },
]);
