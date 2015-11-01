<?php

namespace newsletters\Repositories;

use Prettus\Repository\Eloquent\BaseRepository;
use Prettus\Repository\Criteria\RequestCriteria;
use newsletters\Repositories\ComplaintRepository;
use newsletters\Entities\Complaint;

/**
 * Class ComplaintRepositoryEloquent
 * @package namespace newsletters\Repositories;
 */
class ComplaintRepositoryEloquent extends BaseRepository implements ComplaintRepository
{
    /**
     * Specify Model class name
     *
     * @return string
     */
    public function model()
    {
        return Complaint::class;
    }

    /**
     * Boot up the repository, pushing criteria
     */
    public function boot()
    {
        $this->pushCriteria(app(RequestCriteria::class));
    }
}
