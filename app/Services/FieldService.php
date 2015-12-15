<?php
/**
 * Created by PhpStorm.
 * User: filip
 * Date: 17.8.15
 * Time: 20:00
 */

namespace newsletters\Services;

use Illuminate\Database\QueryException;
use Illuminate\Support\Collection;
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
                $key = $fields->search(function ($field) use ($fieldData) {
                    return strtolower($field->name) === strtolower($fieldData['name']);
                });

                if ($key !== false) {
                    $subscriber->fields()->attach($fields[$key]->id, ['value' => $fieldData['value']]);
                }
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
     * @param int $perPage
     * @param array $with
     * @param array $columns
     * @return mixed
     */
    public function findFieldsByListIdPaginated($listId, $perPage = 10, $with = [], $columns = ['*'])
    {
        return $this->fieldRepository
            ->with($with)
            ->scopeQuery(function ($q) use ($listId) {
                return $q->where('list_id', $listId);
            })->paginate($perPage, $columns);
    }

    /**
     * Find fields by list id
     *
     * @param $listId
     * @param array $with
     * @param array $columns
     * @return mixed
     */
    public function findFieldsByListId($listId, $with = [], $columns = ['*'])
    {
        return $this->fieldRepository
            ->with($with)
            ->scopeQuery(function ($q) use ($listId) {
                return $q->where('list_id', $listId);
            })->all($columns);
    }

    /**
     * Detach subscriber fields by list id
     *
     * @param $listId
     * @param Subscriber $subscriber
     * @return int
     */
    public function detachSubscriberByListId($listId, Subscriber $subscriber)
    {
        $fields = $this->findFieldsByListId($listId, [], ['id'])->toArray();

        return $subscriber->fields()->detach(array_flatten($fields));
    }

    /**
     * Creates an array that is used for the header in the subscribers export file
     *
     * @param $listId
     * @return array
     */
    public function makeHeaderForExportFileByListId($listId)
    {
        return $this->findFieldsByListId($listId)
            ->map(function ($field) {
                return $field->name;
            })
            ->prepend('email')
            ->prepend('name')
            ->toArray();
    }

    /**
     * Find all fields
     *
     * @return mixed
     */
    public function findAllFields()
    {
        return $this->fieldRepository->all();
    }

    /** Find all fields paginated
     *
     * @param int $perPage
     * @return mixed
     */
    public function findAllFieldsPaginated($perPage = 10)
    {
        return $this->fieldRepository->paginate($perPage);
    }

    /**
     * Find a field by id
     *
     * @param $id
     * @return mixed|null
     */
    public function findField($id)
    {
        return $this->fieldRepository->find($id);
    }

    /**
     * Create field
     *
     * @param array $data
     * @return mixed|null
     */
    public function createField(array $data)
    {
        return $this->fieldRepository->create($data);
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
        return $this->fieldRepository->update($data, $id);
    }

    /**
     * Delete a field by its id
     *
     * @param $id
     * @return bool|int
     */
    public function deleteField($id)
    {
        return $this->fieldRepository->delete($id);
    }
}
