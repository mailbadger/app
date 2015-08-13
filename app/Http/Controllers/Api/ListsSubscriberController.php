<?php

namespace newsletters\Http\Controllers\Api;

use Illuminate\Http\Request;

use newsletters\Http\Requests;
use newsletters\Http\Controllers\Controller;
use newsletters\Http\Requests\ImportSubscribersRequest;
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
     * @param ImportSubscribersRequest $request
     * @param $listId
     * @return Response
     */
    public function store(ImportSubscribersRequest $request, $listId)
    {
        //
    }

    /**
     * Display the specified resource.
     *
     * @param  int  $id
     * @return Response
     */
    public function show($listId, $id)
    {
        //
    }

    /**
     * Update the specified resource in storage.
     *
     * @param  Request  $request
     * @param  int  $id
     * @return Response
     */
    public function update(Request $request, $listId, $id)
    {
        //
    }

    /**
     * Remove the specified resource from storage.
     *
     * @param  int  $id
     * @return Response
     */
    public function destroy($listId, $id)
    {
        //
    }
}
