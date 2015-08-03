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
use Illuminate\Support\Facades\Log;
use newsletters\Repositories\ListsRepository;

class ListsService
{

    /**
     * @var ListsRepository
     */
    private $listsRepository;

    public function __construct(ListsRepository $repository)
    {
        $this->listsRepository = $repository;
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
}