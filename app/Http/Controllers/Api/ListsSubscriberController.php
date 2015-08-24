<?php

namespace newsletters\Http\Controllers\Api;

use Illuminate\Http\Request;
use newsletters\Http\Controllers\Controller;
use newsletters\Http\Requests;
use newsletters\Http\Requests\ImportSubscribersRequest;
use newsletters\Services\FieldService;
use newsletters\Services\FileService;
use newsletters\Services\ListsService;

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
    public function index(Request $request, $listId)
    {
        $subscribers = $this->service->findAllSubscribersByListId($listId, $request->has('paginate'), 10);

        return response()->json($subscribers, 200);
    }

    /**
     * Store a newly created resource in storage.
     *
     * @return Response
     */
    public function store()
    {
        //
    }

    /**
     * Display the specified resource.
     *
     * @param  int $id
     * @return Response
     */
    public function show($listId, $id)
    {
        //
    }

    /**
     * Update the specified resource in storage.
     *
     * @param  Request $request
     * @param  int $id
     * @return Response
     */
    public function update(Request $request, $listId, $id)
    {
        //
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
        if ($this->service->detachSubscriber($listId, $id)) {
            return response()->json(['status' => 200, 'message' => 'The specified resource has been deleted.'],
                200);
        }

        return response()->json(['status' => 422, 'campaign' => ['The specified resource could not be deleted.']],
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
        if (!$subscribers->isEmpty()) {
            return response()->json(['status' => 200, 'subscribers' => 'The specified resources have been created.'],
                200);
        }

        return response()->json(['status' => 412, 'subscribers' => ['The specified resource could not be created.']],
            412);
    }

    public function export(
        Request $request,
        FileService $fileService,
        FieldService $fieldService,
        $listId
    ) {
        $excel = $this->service->exportSubscribers($listId, $fileService, $fieldService);

        header('Content-Type: text/csv');
        header('Content-Disposition: attachment;filename="Subscribers_'.date('dMy').'.csv"');
        header('Cache-Control: max-age=0');

        $writer = $fileService->createWriter($excel);
        $writer->save('php://output');

        exit;
    }
}
