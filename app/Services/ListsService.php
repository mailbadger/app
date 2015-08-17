<?php
/**
 * Created by PhpStorm.
 * User: filip
 * Date: 3.8.15
 * Time: 17:41
 */

namespace newsletters\Services;


use Exception;
use Illuminate\Database\Eloquent\ModelNotFoundException;
use Illuminate\Database\QueryException;
use Illuminate\Support\Collection;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;
use newsletters\Repositories\ListsRepository;
use newsletters\Repositories\SubscriberRepository;

class ListsService
{

    /**
     * @var ListsRepository
     */
    private $listsRepository;

    /**
     * @var SubscriberRepository
     */
    private $subscriberRepository;

    public function __construct(ListsRepository $listsRepository, SubscriberRepository $subscriberRepository)
    {
        $this->listsRepository = $listsRepository;
        $this->subscriberRepository = $subscriberRepository;
    }

    /**
     * Find all subscribers lists
     *
     * @param bool|false $paginate
     * @param int $perPage
     * @return mixed
     */
    public function findAllLists($paginate = false, $perPage = 10)
    {
        if ($paginate) {
            return $this->listsRepository->paginate($perPage);
        }

        return $this->listsRepository->all();
    }

    /**
     * Find all subscribers on a list
     * @param $listId
     * @param bool|false $paginate
     * @param int $perPage
     * @return mixed
     */
    public function findAllSubscribersByListId($listId, $paginate = false, $perPage = 10)
    {
        $subscribers = $this->listsRepository->find($listId)->subscribers();

        return (!empty($paginate)) ? $subscribers->paginate($perPage) : $subscribers->all();
    }

    /**
     * Find a list by id
     *
     * @param $id
     * @return mixed|null
     */
    public function findList($id)
    {
        try {
            return $this->listsRepository->find($id);
        } catch (ModelNotFoundException $e) {
            return null;
        }
    }

    /**
     * Create list
     *
     * @param array $data
     * @return mixed|null
     */
    public function createList(array $data)
    {
        try {
            return $this->listsRepository->create($data);
        } catch (Exception $e) {
            Log::error($e->getMessage());

            return null;
        }
    }

    /**
     * Update list by id
     *
     * @param array $data
     * @param $id
     * @return mixed|null
     */
    public function updateList(array $data, $id)
    {
        try {
            return $this->listsRepository->update($data, $id);
        } catch (Exception $e) {
            Log::error($e->getMessage());

            return null;
        }
    }

    /**
     * Delete a list by its id
     *
     * @param $id
     * @return bool|int
     */
    public function deleteList($id)
    {
        try {
            return $this->listsRepository->delete($id);
        } catch (ModelNotFoundException $e) {

            return false;
        }
    }

    /**
     * @param $file
     * @param $listId
     * @param FileService $fileService
     * @param FieldService $fieldService
     * @return Collection
     */
    public function createSubscribers($file, $listId, FileService $fileService, FieldService $fieldService)
    {
        return $fileService->importSubscribersFromFile($file)
            ->map(function ($data) use ($listId, $fieldService) {
                return DB::transaction(function () use ($data, $listId, $fieldService) {
                    try {
                        $subscriber = $this->subscriberRepository->create($data['subscriber']);
                        $fieldService->createSubscriberFields($subscriber, $data['custom_fields'], $listId);
                        $subscriber->lists()->attach($listId);

                        return $subscriber;
                    } catch (QueryException $e) {
                        Log::error($e->getMessage() . '\nLine: ' . $e->getLine() . '\nStack trace: ' . $e->getTraceAsString());

                        return null;
                    }
                });
            })
            ->reject(function ($s) {
                return empty($s);
            });
    }
}