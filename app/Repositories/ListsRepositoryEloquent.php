<?php

namespace newsletters\Repositories;

use Prettus\Repository\Eloquent\BaseRepository;
use Prettus\Repository\Criteria\RequestCriteria;
use newsletters\Entities\Lists;

/**
 * Class ListsRepositoryEloquent
 * @package namespace newsletters\Repositories;
 */
class ListsRepositoryEloquent extends BaseRepository implements ListsRepository
{
    /**
     * Specify Model class name
     *
     * @return string
     */
    public function model()
    {
        return Lists::class;
    }

    /**
     * Boot up the repository, pushing criteria
     */
    public function boot()
    {
        $this->pushCriteria( app(RequestCriteria::class) );
    }
}