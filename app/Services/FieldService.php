<?php
/**
 * Created by PhpStorm.
 * User: filip
 * Date: 17.8.15
 * Time: 20:00
 */

namespace newsletters\Services;


use Exception;
use Illuminate\Database\Eloquent\ModelNotFoundException;
use Illuminate\Database\QueryException;
use Illuminate\Support\Collection;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;
use newsletters\Entities\Subscriber;
use newsletters\Repositories\FieldRepository;

class FieldService
{

    /**
     * @var FieldRepository
     */
    private $fieldRepository;

    public function __construct(FieldRepository $fieldRepository)
    {
        $this->fieldRepository = $fieldRepository;
    }

    /**
     * Create multiple fields for subscriber
     *
     * @param Subscriber $subscriber
     * @param array $data
     * @param $listId
     * @return Collection
     */
    public function attachFieldsToSubscriber(Subscriber $subscriber, array $data, $listId)
    {
        try {
            $fields = $this->findFieldsByListId($listId);

            foreach ($data as $fieldData) {
                DB::transaction(function () use ($fieldData, $listId, $subscriber, $fields) {
                    $key = $fields->search(function ($field) use ($fieldData) {
                        return strtolower($field->name) === strtolower($fieldData['name']);
                    });

                    if ($key !== false) {
                        $subscriber->fields()->attach($fields[$key]->id, ['value' => $fieldData['value']]);
                    }
                });
            }

            return true;
        } catch (QueryException $e) {
            Log::error($e->getMessage() . '\nLine: ' . $e->getLine() . '\nStack trace: ' . $e->getTraceAsString());

            return false;
        }
    }

    /**
     * Find a field by its name and by list id
     *
     * @param $name
     * @param $listId
     * @param array $with
     * @param array $columns
     * @return mixed
     */
    public function findFieldByNameAndListId($name, $listId, $with = [], $columns = ['*'])
    {
        return $this->fieldRepository
            ->with($with)
            ->findWhere(['name' => $name, 'list_id' => $listId], $columns)
            ->first();
    }

    /**
     * Find fields by list id
     *
     * @param $listId
     * @param bool $paginate
     * @param int $perPage
     * @param array $with
     * @param array $columns
     * @return mixed
     */
    public function findFieldsByListId($listId, $paginate = false, $perPage = 10, $with = [], $columns = ['*'])
    {
        $query = $this->fieldRepository
            ->with($with)
            ->scopeQuery(function ($q) use ($listId) {
                return $q->where('list_id', $listId);
            });

        return (!empty($paginate)) ? $query->paginate($perPage, $columns) : $query->all($columns);
    }

    /**
     * Find all fields
     *
     * @param bool $paginate
     * @param int $perPage
     * @return mixed
     */
    public function findAllFields($paginate = false, $perPage = 10)
    {
        return (!empty($paginate)) ? $this->fieldRepository->paginate($perPage) : $this->fieldRepository->all();
    }

    /**
     * Find a field by id
     *
     * @param $id
     * @return mixed|null
     */
    public function findField($id)
    {
        try {
            return $this->fieldRepository->find($id);
        } catch (ModelNotFoundException $e) {
            return null;
        }
    }

    /**
     * Create field
     *
     * @param array $data
     * @return mixed|null
     */
    public function createField(array $data)
    {
        try {
            return $this->fieldRepository->create($data);
        } catch (Exception $e) {
            Log::error($e->getMessage());

            return null;
        }
    }

    /**
     * Update field by id
     *
     * @param array $data
     * @param $id
     * @return mixed|null
     */
    public function updateField(array $data, $id)
    {
        try {
            return $this->fieldRepository->update($data, $id);
        } catch (Exception $e) {
            Log::error($e->getMessage());

            return null;
        }
    }

    /**
     * Delete a field by its id
     *
     * @param $id
     * @return bool|int
     */
    public function deleteField($id)
    {
        try {
            return $this->fieldRepository->delete($id);
        } catch (ModelNotFoundException $e) {

            return false;
        }
    }
}