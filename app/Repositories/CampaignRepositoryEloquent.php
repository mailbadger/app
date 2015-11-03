<?php

namespace newsletters\Repositories;

use newsletters\Entities\Campaign;
use Prettus\Repository\Criteria\RequestCriteria;
use Prettus\Repository\Eloquent\BaseRepository;

/**
 * Class CampaignRepositoryEloquent
 * @package namespace newsletters\Repositories;
 */
class CampaignRepositoryEloquent extends BaseRepository implements CampaignRepository
{
    /**
     * @var array
     */
    protected $fieldSearchable = [
        'name',
        'status',
    ];

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
        $this->pushCriteria(app(RequestCriteria::class));
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
