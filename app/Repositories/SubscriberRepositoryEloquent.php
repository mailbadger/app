<?php

namespace newsletters\Repositories;

use Prettus\Repository\Eloquent\BaseRepository;
use Prettus\Repository\Criteria\RequestCriteria;
use newsletters\Entities\Subscriber;

/**
 * Class SubscriberRepositoryEloquent
 * @package namespace newsletters\Repositories;
 */
class SubscriberRepositoryEloquent extends BaseRepository implements SubscriberRepository
{
    /**
     * Specify Model class name
     *
     * @return string
     */
    public function model()
    {
        return Subscriber::class;
    }

    /**
     * Boot up the repository, pushing criteria
     */
    public function boot()
    {
        $this->pushCriteria( app(RequestCriteria::class) );
    }
}