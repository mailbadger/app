<?php
/**
 * Created by PhpStorm.
 * User: filip
 * Date: 3.8.15
 * Time: 17:41
 */

namespace newsletters\Services;

use Exception;
use Illuminate\Support\Collection;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;
use newsletters\Entities\Lists;
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
     * Find a list by id
     *
     * @param $id
     * @return mixed|null
     */
    public function findList($id)
    {
        return $this->listsRepository->find($id);
    }

    /**
     * Create list
     *
     * @param array $data
     * @return mixed|null
     */
    public function createList(array $data)
    {
        return $this->listsRepository->create($data);
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
        return $this->listsRepository->update($data, $id);
    }

    /**
     * Delete a list by its id
     *
     * @param $id
     * @return bool|int
     */
    public function deleteList($id)
    {
        return $this->listsRepository->delete($id);
    }

    /**
     * Create subscribers from a imported file
     *
     * @param $file
     * @param $listId
     * @param FileService $fileService
     * @param FieldService $fieldService
     * @return Collection
     * @throws Exception
     */
    public function createSubscribers($file, $listId, FileService $fileService, FieldService $fieldService)
    {
        $list = $this->findList($listId);

        $totalSubscribers = $fileService->importSubscribers($file)
            ->map(function ($data) use ($list, $fieldService) {
                return DB::transaction(function () use ($data, $list, $fieldService) {
                    try {
                        $subscriber = $this->subscriberRepository->findByField('email',
                            $data['subscriber']['email'])->first();

                        if (empty($subscriber)) {
                            $subscriber = $this->subscriberRepository->create($data['subscriber']);
                        }

                        $fieldService->attachFieldsToSubscriber($subscriber, $data['custom_fields'], $list->id);
                        $this->attachSubscriber($list, $subscriber->id);

                        return $subscriber;
                    } catch (Exception $e) {
                        Log::error($e->getMessage() . '\nLine: ' . $e->getLine() . '\nStack trace: ' . $e->getTraceAsString());

                        return null;
                    }
                });
            })
            ->reject(function ($s) {
                return empty($s);
            })
            ->count();

        $total = $list->total_subscribers + $totalSubscribers;
        $this->updateTotalListSubscribers($list, $total);

        return $totalSubscribers;
    }

    /**
     * Mass detach subscribers from a list by importing a file with emails
     *
     * @param $file
     * @param $listId
     * @param FileService $fileService
     * @param FieldService $fieldService
     * @return int
     */
    public function deleteSubscribers($file, $listId, FileService $fileService, FieldService $fieldService)
    {
        $list = $this->findList($listId);

        $subscribers = $fileService->importSubscribers($file)
            ->map(function ($data) use ($list, $fieldService) {
                $subscriber = $this->subscriberRepository->findByField('email', $data['subscriber']['email'],
                    ['id'])->first();

                if (!empty($subscriber)) {
                    $fieldService->detachSubscriberByListId($list->id, $subscriber);
                    $this->detachSubscriber($list, $subscriber->id);

                    return true;
                }

                return false;
            })
            ->reject(function ($s) {
                return empty($s);
            })
            ->count();

        $total = $list->total_subscribers - $subscribers;
        $this->updateTotalListSubscribers($list, $total);

        return $subscribers;
    }

    /**
     * Exports subscribers to an excel file
     *
     * @param $listId
     * @param FileService $fileService
     * @param FieldService $fieldService
     * @return \PHPExcel
     */
    public function exportSubscribers($listId, FileService $fileService, FieldService $fieldService)
    {
        $header = $fieldService->findFieldsByListId($listId)
            ->map(function ($field) {
                return $field->name;
            })
            ->prepend('email')
            ->prepend('name')
            ->toArray();

        $subscribers = $this->subscriberRepository
            ->with('fields')
            ->scopeQuery(function ($q) use ($listId) {
                return $q->whereHas('lists', function ($q) use ($listId) {
                    return $q->where('list_id', $listId);
                });
            })
            ->all()
            ->map(function ($sub) {
                $sub->fields->each(function ($field) use ($sub) {
                    $name = $field->name;
                    $sub->$name = $field->pivot->value;
                });

                unset($sub->fields);

                return $sub;
            })
            ->toArray();

        $excelObj = $fileService->createExcelFile('subscribers');

        return $fileService->exportData($subscribers, $header, $excelObj);
    }

    /**
     * Attaches a subscriber to a list
     *
     * @param Lists $list
     * @param $id
     * @return void
     */
    public function attachSubscriber(Lists $list, $id)
    {
        $list->subscribers()->attach($id);
    }

    /**
     * Detaches a subscriber from a list
     *
     * @param Lists $list
     * @param $id
     * @return bool
     */
    public function detachSubscriber(Lists $list, $id)
    {
        return $list->subscribers()->detach($id);
    }

    /**
     * Update total subscribers to list
     *
     * @param Lists $list
     * @param $total
     * @return void 
     */
    public function updateTotalListSubscribers(Lists $list, $total)
    {
        $list->total_subscribers = (!is_numeric($total) || $total < 0) ? 0 : $total;
        $list->save();
    }
}
