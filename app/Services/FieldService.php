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
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;
use newsletters\Entities\Subscriber;
use newsletters\Repositories\FieldRepository;

class FieldService
{

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
    public function createSubscriberFields(Subscriber $subscriber, array $data, $listId)
    {
        $fields = new Collection();

        try {
            foreach ($data as $fieldData) {
                $field = DB::transaction(function () use ($fieldData, $listId, $subscriber) {
                    $field = $this->createSubscriberField($fieldData['name'], $listId);
                    $subscriber->fields()->attach($field->id, ['value' => $fieldData['value']]);

                    return $field;
                });

                $fields->push($field);
            }
        } catch (QueryException $e) {
            Log::error($e->getMessage() . '\nLine: ' . $e->getLine() . '\nStack trace: ' . $e->getTraceAsString());
        }

        return $fields;
    }

    /**
     * Create and associate field with subscriber
     *
     * @param $name
     * @param $listId
     * @return mixed|null
     */
    public function createSubscriberField($name, $listId)
    {
        try {
            $field = $this->findFieldByNameAndListId($name, $listId);

            if (empty($field)) {
                $field = $this->fieldRepository->create(['name' => $name, 'list_id' => $listId]);
            }

            return $field;
        } catch (QueryException $e) {
            Log::error($e->getMessage() . '\nLine: ' . $e->getLine() . '\nStack trace: ' . $e->getTraceAsString());

            return null;
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
}