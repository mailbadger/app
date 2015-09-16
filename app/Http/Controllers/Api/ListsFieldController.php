<?php

namespace newsletters\Http\Controllers\Api;

use Illuminate\Http\Request;
use newsletters\Http\Controllers\Controller;
use newsletters\Http\Requests;
use newsletters\Services\FieldService;

class ListsFieldController extends Controller
{

    /**
     * @var FieldService
     */
    private $service;

    public function __construct(FieldService $service)
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
        $subscribers = $this->service->findFieldsByListId($listId, $request->has('paginate'), 10);

        return response()->json($subscribers, 200);
    }

    /**
     * Store a newly created resource in storage.
     *
     * @param  Request $request
     * @param $listId
     * @return Response
     */
    public function store(Request $request, $listId)
    {
        $data = $request->all();
        $data['list_id'] = $listId;
        $field = $this->service->createField($data);

        if (isset($field)) {
            return response()->json(['field' => $field->id], 200);
        }

        return response()->json(['message' => ['The specified resource could not be created.']], 412);
    }

    /**
     * Display the specified resource.
     *
     * @param  int $id
     * @return Response
     */
    public function show($listId, $id)
    {
        $field = $this->service->findField($id);

        if (isset($field)) {
            return response()->json($field, 200);
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
    public function update(Request $request, $listId, $id)
    {
        $field = $this->service->updateField($request->all(), $id);
        if (isset($field)) {
            return response()->json(['field' => $field->id], 200);
        }

        return response()->json(['message' => ['The specified resource could not be updated.']], 412);
    }

    /**
     * Remove the specified resource from storage.
     *
     * @param  int $id
     * @return Response
     */
    public function destroy($listId, $id)
    {
        if ($this->service->deleteField($id)) {
            return response()->json(['message' => 'The specified resource has been deleted.'], 200);
        }

        return response()->json(['message' => ['The specified resource could not be deleted.']], 422);
    }
}
