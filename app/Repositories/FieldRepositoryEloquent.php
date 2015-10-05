<?php

namespace newsletters\Repositories;

use newsletters\Entities\Field;
use Prettus\Repository\Criteria\RequestCriteria;
use Prettus\Repository\Eloquent\BaseRepository;

/**
 * Class FieldRepositoryEloquent
 * @package namespace newsletters\Repositories;
 */
class FieldRepositoryEloquent extends BaseRepository implements FieldRepository
{
    /**
     * Specify Model class name
     *
     * @return string
     */
    public function model()
    {
        return Field::class;
    }

    /**
     * Boot up the repository, pushing criteria
     */
    public function boot()
    {
        $this->pushCriteria(app(RequestCriteria::class));
    }
}
