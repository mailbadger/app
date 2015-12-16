<?php

namespace newsletters\Http\Controllers\Api;

use Illuminate\Http\Request;
use newsletters\Http\Controllers\Controller;
use newsletters\Http\Requests\ImportSubscribersRequest;
use newsletters\Http\Requests\MassDeleteSubscribersRequest;
use newsletters\Services\FieldService;
use newsletters\Services\FileService;
use newsletters\Services\ListsService;
use newsletters\Services\SubscriberService;

class ListsSubscriberController extends Controller
{
    /**
     * @var ListsService
     */
    private $service;

    public function __construct(ListsService $service)
    {
        $this->middleware('auth.basic');

        $this->service = $service;
    }

    /**
     * Display a listing of the resource.
     *
     * @param Request $request
     * @param $listId
     * @return Response
     */
    public function index(Request $request, $listId, SubscriberService $subscriberService)
    {
        $perPage = ($request->has('per_page')) ? $request->input('per_page') : 10;

        if ($request->has('paginate')) {
            $subscribers = $subscriberService->findAllSubscribersByListIdPaginated($listId, $perPage);
        } else {
            $subscribers = $subscriberService->findAllSubscribersByListId($listId);
        }

        return response()->json($subscribers, 200);
    }

    /**
     * Store a newly created resource in storage.
     *
     * @return Response
     */
    public function store(StoreSubscriberRequest $request, $listId, SubscriberService $subscriberService)
    {
        $subscriber = $subscriberService->findSubscriberByEmail($request->input('email'));

        if (empty($subscriber)) {
            $subscriber = $subscriberService->createSubscriber($request->all());
        }

        if (!empty($subscriber)) {
            $list = $this->service->findList($listId);
            $this->service->attachSubscriber($list, $subscriber->id);

            return response()->json(['subscriber' => $subscriber->id], 200);
        }

        return response()->json(['errors' => ['The specified resource could not be created.']], 412);
    }

    /**
     * Display the specified resource.
     *
     * @param int $id
     * @return Response
     */
    public function show($listId, $id, SubscriberService $subscriberService)
    {
        $subscriber = $subscriberService->findSubscriber($id);

        if (isset($subscriber)) {
            return response()->json($subscriber, 200);
        }

        return response()->json(['errors' => ['The specified resource does not exist.']], 404);
    }

    /**
     * Update the specified resource in storage.
     *
     * @param  Request $request
     * @param $listId
     * @param  int $id
     * @return Response
     */
    public function update(Request $request, $listId, $id, SubscriberService $subscriberService)
    {
        $subscriber = $subscriberService->updateSubscriber($request->all(), $id);

        if (isset($subscriber)) {
            return response()->json(['subscriber' => $subscriber->id], 200);
        }

        return response()->json(['errors' => ['The specified resource could not be updated.']], 412);
    }

    /**
     * Remove the specified resource from storage.
     *
     * @param $listId
     * @param  int $id
     * @return Response
     */
    public function destroy($listId, $id)
    {
        $list = $this->service->findList($listId);

        if ($this->service->detachSubscriber($list, $id)) {
            $this->service->updateTotalListSubscribers($list, $list->total_subscribers - 1);
            return response()->json(['message' => 'The specified resource has been deleted.'],
                200);
        }

        return response()->json(['message' => ['The specified resource could not be deleted.']],
            422);
    }

    /**
     * Import subscribers from a file
     *
     * @param ImportSubscribersRequest $request
     * @param FileService $fileService
     * @param FieldService $fieldService
     * @param $listId
     * @return \Illuminate\Http\JsonResponse
     */
    public function import(
        ImportSubscribersRequest $request,
        FileService $fileService,
        FieldService $fieldService,
        $listId
    ) {
        $subscribers = $this->service->createSubscribers($request->file('subscribers'), $listId, $fileService,
            $fieldService);
        if (!empty($subscribers)) {
            return response()->json(['message' => 'The specified resources have been created.'],
                200);
        }

        return response()->json(['message' => ['The specified resource could not be created.']],
            412);
    }

    /**
     * Export subscribers to csv file
     *
     * @param FileService $fileService
     * @param FieldService $fieldService
     * @param $listId
     */
    public function export(
        FileService $fileService,
        FieldService $fieldService,
        $listId
    ) {
        $excel = $this->service->exportSubscribers($listId, $fileService, $fieldService);

        header('Content-Type: text/csv');
        header('Content-Disposition: attachment;filename="Subscribers_' . date('dMy') . '.csv"');
        header('Cache-Control: max-age=0');

        $writer = $fileService->createWriter($excel);
        $writer->save('php://output');

        exit;
    }

    /**
     * Mass delete subscribers from a file
     *
     * @param MassDeleteSubscribersRequest $request
     * @param FileService $fileService
     * @param FieldService $fieldService
     * @param $listId
     * @return \Illuminate\Http\JsonResponse
     */
    public function massDelete(
        MassDeleteSubscribersRequest $request,
        FileService $fileService,
        FieldService $fieldService,
        $listId
    ) {
        $count = $this->service->deleteSubscribers($request->file('subscribers'), $listId, $fileService, $fieldService);

        if (!empty($count)) {
            return response()->json(['message' => 'The specified resource have been deleted.'],
                200);
        }

        return response()->json(['message' => ['The specified resource could not be deleted.']],
            422);
    }
}
