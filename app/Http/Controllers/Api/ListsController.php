<?php

namespace newsletters\Http\Controllers\Api;

use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;
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
        $perPage = ($request->has('per_page')) ? $request->input('per_page') : 10;

        if($request->has('paginate')) {
            $lists = $this->service->findAllListsPaginated($perPage);
        } else {
            $lists = $this->service->findAllLists();
        }

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

        if(isset($list)) {
            return response()->json(['list' => $list->id], 200);
        }

        return response()->json(['message' => ['The specified resource could not be created.']], 422);
    }

    /**
     * Display the specified resource.
     *
     * @param  int $id
     * @return Response
     */
    public function show($id)
    {
        $list = $this->service->findList($id);

        if(isset($list)) {
            return response()->json($list, 200);
        }

        return response()->json(['message' => 'The specified resource does not exist.'], 404);
    }

    /**
     * Update the specified resource in storage.
     *
     * @param  Request $request
     * @param  int $id
     * @return Response
     */
    public function update(Request $request, $id)
    {
        $list = $this->service->updateList($request->all(), $id);

        if(isset($list)) {
            return response()->json(['list' => $list->id], 200);
        }

        return response()->json(['message' => ['The specified resource could not be updated.']], 422);
    }

    /**
     * Remove the specified resource from storage.
     *
     * @param  int $id
     * @return Response
     */
    public function destroy($id)
    {
        if ($this->service->deleteList($id)) {
            return response()->json(['message' => ['The specified resource has been deleted.']], 200);
        }

        return response()->json(['message' => ['The specified resource could not be deleted.']], 422);
    }
}
