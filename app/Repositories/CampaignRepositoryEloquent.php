<?php

namespace newsletters\Repositories;

use Prettus\Repository\Eloquent\BaseRepository;
use Prettus\Repository\Criteria\RequestCriteria;
use newsletters\Entities\Campaign;

/**
 * Class CampaignRepositoryEloquent
 * @package namespace newsletters\Repositories;
 */
class CampaignRepositoryEloquent extends BaseRepository implements CampaignRepository
{
    /**
     * Specify Model class name
     *
     * @return string
     */
    public function model()
    {
        return Campaign::class;
    }

    /**
     * Boot up the repository, pushing criteria
     */
    public function boot()
    {
        $this->pushCriteria( app(RequestCriteria::class) );
    }

    /**
     * Load trashed entities
     *
     * @return $this
     */
    public function withTrashed()
    {
        $this->model = $this->model->withTrashed();
        return $this;
    }
}