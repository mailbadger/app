<?php

namespace newsletters\Repositories;

use newsletters\Entities\Template;
use Prettus\Repository\Criteria\RequestCriteria;
use Prettus\Repository\Eloquent\BaseRepository;

/**
 * Class TemplateRepositoryEloquent
 * @package namespace newsletters\Repositories;
 */
class TemplateRepositoryEloquent extends BaseRepository implements TemplateRepository
{
    /**
     * Specify Model class name
     *
     * @return string
     */
    public function model()
    {
        return Template::class;
    }

    /**
     * Boot up the repository, pushing criteria
     */
    public function boot()
    {
        $this->pushCriteria(app(RequestCriteria::class));
    }
}
