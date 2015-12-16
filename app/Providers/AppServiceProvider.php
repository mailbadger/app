<?php

namespace newsletters\Providers;

use Illuminate\Contracts\Validation\Factory;
use Illuminate\Support\ServiceProvider;

class AppServiceProvider extends ServiceProvider
{
    /**
     * Bootstrap any application services.
     *
     * @return void
     */
    public function boot()
    {
        $this->registerValidationRules($this->app['validator']);

        $this->bindRepositories();
    }

    /**
     * Register custom validation rules here with the Factory validator
     *
     * @param Factory $validator
     */
    public function registerValidationRules(Factory $validator)
    {
        $validator->extend('check_fields', 'newsletters\Validators\ListValidator@validateCheckFields');

        $validator->extend('array_of_emails', 'newsletters\Validators\CampaignValidator@validateArrayOfEmails');
    }

    public function bindRepositories()
    {
        $this->app->bind('newsletters\Repositories\CampaignRepository',
            'newsletters\Repositories\CampaignRepositoryEloquent');

        $this->app->bind('newsletters\Repositories\TemplateRepository',
            'newsletters\Repositories\TemplateRepositoryEloquent');

        $this->app->bind('newsletters\Repositories\ListsRepository',
            'newsletters\Repositories\ListsRepositoryEloquent');

        $this->app->bind('newsletters\Repositories\TagRepository',
            'newsletters\Repositories\TagRepositoryEloquent');

        $this->app->bind('newsletters\Repositories\SubscriberRepository',
            'newsletters\Repositories\SubscriberRepositoryEloquent');

        $this->app->bind('newsletters\Repositories\FieldRepository',
            'newsletters\Repositories\FieldRepositoryEloquent');

        $this->app->bind('newsletters\Repositories\SentEmailRepository',
            'newsletters\Repositories\SentEmailRepositoryEloquent');

        $this->app->bind('newsletters\Repositories\UserRepository',
            'newsletters\Repositories\UserRepositoryEloquent');

        $this->app->bind('newsletters\Repositories\BounceRepository',
            'newsletters\Repositories\BounceRepositoryEloquent');

        $this->app->bind('newsletters\Repositories\ComplaintRepository',
            'newsletters\Repositories\ComplaintRepositoryEloquent');
    }

    /**
     * Register any application services.
     *
     * @return void
     */
    public function register()
    {
        //
    }
}
