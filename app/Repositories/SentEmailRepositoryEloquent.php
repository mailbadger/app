<?php

namespace newsletters\Repositories;

use newsletters\Entities\SentEmail;
use Prettus\Repository\Criteria\RequestCriteria;
use Prettus\Repository\Eloquent\BaseRepository;

/**
 * Class SentEmailRepositoryEloquent
 * @package namespace newsletters\Repositories;
 */
class SentEmailRepositoryEloquent extends BaseRepository implements SentEmailRepository
{
    /**
     * Specify Model class name
     *
     * @return string
     */
    public function model()
    {
        return SentEmail::class;
    }

    /**
     * Boot up the repository, pushing criteria
     */
    public function boot()
    {
        $this->pushCriteria(app(RequestCriteria::class));
    }
}