<?php

namespace newsletters\Repositories;

use newsletters\Entities\Subscriber;
use Prettus\Repository\Criteria\RequestCriteria;
use Prettus\Repository\Eloquent\BaseRepository;
use Closure;

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
        $this->pushCriteria(app(RequestCriteria::class));
    }

    /**
     * Returns chunked data
     * @param int $chunks
     * @param Closure $closure
     * @return mixed
     */
    public function chunk($chunks, Closure $closure)
    {
        $this->applyCriteria();
        $this->applyScope();

        $results = $this->model->chunk($chunks, $closure);
  
        $this->resetModel();

        return $this->parserResult($results);
    }
}
