<?php

namespace newsletters\Services;

use newsletters\Repositories\SubscriberRepository;
use Closure;

class SubscriberService
{
    /**
     * @var SubscriberRepository
     */    
    protected $subscriberRepository;

    public function __construct(SubscriberRepository $subscriberRepository)
    {
        $this->subscriberRepository = $subscriberRepository;
    }

    /**
     * Find a subscriber by his email address
     * @param $email
     * @return mixed
     */
    public function findSubscriberByEmail($email)
    {
        return $this->subscriberRepository->findByField('email', $email);
    }

    /**
     * Find all subscribers on a list
     *
     * @param $listId
     * @return mixed
     */
    public function findAllSubscribersByListId($listId, $paginate = false, $perPage = 10)
    {
        return $this->subscriberRepository->scopeQuery(function ($q) use ($listId) {
            return $q->whereHas('lists', function ($q) use ($listId) {
                return $q->where('list_id', $listId);
            });
        })->all();
    }

    /**
     * Find all subscribers on a list paginated
     * 
     * @param $listId
     * @param $perPage
     * @return mixed
     */
    public function findAllSubscribersByListIdPaginated($listId, $perPage = 10)
    {
       return $this->subscriberRepository->scopeQuery(function ($q) use ($listId) {
            return $q->whereHas('lists', function ($q) use ($listId) {
                return $q->where('list_id', $listId);
            });
        })->paginate($perPage);
    }

    /**
     * Find all subscribers by list ids. Returns chunks of data for processing
     *
     * @param array $listIds
     * @param int $chunks
     * @param Closure $closure
     * @return mixed
     */
    public function findSubscribersByListIdsByChunks(array $listIds, $chunks = 1000, Closure $closure)
    {
        return $this->subscriberRepository
            ->with('fields')
            ->scopeQuery(function ($q) use ($listIds) {
                return $q->distinct()->whereHas('lists', function ($q) use ($listIds) {
                    return $q->whereIn('list_id', $listIds);
                });
            })
            ->chunk($chunks, $closure); 
    }

    /**
     * Find a subscriber by id
     *
     * @param $id
     * @return mixed|null
     */
    public function findSubscriber($id)
    {
        return $this->subscriberRepository->find($id);
    }


    /**
     * Create subscriber
     *
     * @param array $data
     * @return mixed 
     */
    public function createSubscriber(array $data)
    {
        return $this->subscriberRepository->create($data);
    }
    /**
     * Update subscriber by id
     *
     * @param array $data
     * @param int $id
     * @return int
     */
    public function updateSubscriber(array $data, $id)
    {
        return $this->subscriberRepository->update($data, $id);
    }
}
