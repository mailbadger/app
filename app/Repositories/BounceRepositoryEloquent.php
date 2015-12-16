<?php

namespace newsletters\Repositories;

use Prettus\Repository\Eloquent\BaseRepository;
use Prettus\Repository\Criteria\RequestCriteria;
use newsletters\Repositories\BounceRepository;
use newsletters\Entities\Bounce;

/**
 * Class BounceRepositoryEloquent
 * @package namespace newsletters\Repositories;
 */
class BounceRepositoryEloquent extends BaseRepository implements BounceRepository
{
    /**
     * Specify Model class name
     *
     * @return string
     */
    public function model()
    {
        return Bounce::class;
    }

    /**
     * Boot up the repository, pushing criteria
     */
    public function boot()
    {
        $this->pushCriteria(app(RequestCriteria::class));
    }
}
