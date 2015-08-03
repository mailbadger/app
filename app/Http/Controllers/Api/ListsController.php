<?php

namespace newsletters\Http\Controllers\Api;

use Illuminate\Http\Request;

use newsletters\Http\Requests;
use newsletters\Http\Controllers\Controller;
use newsletters\Http\Requests\StoreListRequest;
use newsletters\Services\ListsService;

class ListsController extends Controller
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
     * @return Response
     */
    public function index(Request $request)
    {
        $lists = $this->service->findAllLists($request->has('paginate'), 10);

        return response()->json($lists, 200);
    }

    /**
     * Store a newly created resource in storage.
     *
     * @param StoreListRequest $request
     * @return Response
     */
    public function store(StoreListRequest $request)
    {
        $list = $this->service->createList($request->all());
        if (isset($list)) {
            return response()->json(['status' => 200, 'list' => $list->id], 200);
        }

        return response()->json(['status' => 412, 'list' => ['The specified resource could not be created.']],
            412);
    }

    /**
     * Display the specified resource.
     *
     * @param  int  $id
     * @return Response
     */
    public function show($id)
    {
        $list = $this->service->findList($id);

        if (isset($list)) {
            return response()->json($list, 200);
        }

        return response()->json(['status' => 404, 'message' => 'The specified resource does not exist.'], 404);
    }

    /**
     * Update the specified resource in storage.
     *
     * @param  Request  $request
     * @param  int  $id
     * @return Response
     */
    public function update(Request $request, $id)
    {
        $list = $this->service->updateList($request->all(), $id);
        if (isset($list)) {
            return response()->json(['status' => 200, 'list' => $list->id], 200);
        }

        return response()->json(['status' => 412, 'list' => ['The specified resource could not be updated.']],
            412);
    }

    /**
     * Remove the specified resource from storage.
     *
     * @param  int  $id
     * @return Response
     */
    public function destroy($id)
    {
        if ($this->service->deleteList($id)) {
            return response()->json(['status' => 200, 'message' => 'The specified resource has been deleted.'],
                200);
        }

        return response()->json(['status' => 422, 'list' => ['The specified resource could not be deleted.']],
            422);
    }
}
