<?php

namespace newsletters\Http\Controllers\Api;

use Illuminate\Http\Request;

use newsletters\Http\Requests;
use newsletters\Http\Controllers\Controller;
use newsletters\Services\FieldService;

class FieldController extends Controller
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
     * @return Response
     */
    public function index(Request $request)
    {
        $fields = $this->service->findAllFields($request->has('paginate'), 10);

        return response()->json($fields, 200);
    }

    /**
     * Display the specified resource.
     *
     * @param  int  $id
     * @return Response
     */
    public function show($id)
    {

        $field = $this->service->findField($id);

        if (isset($field)) {
            return response()->json($field, 200);
        }

        return response()->json(['status' => 404, 'message' => 'The specified resource does not exist.'], 404);
    }

    /**
     * Remove the specified resource from storage.
     *
     * @param  int  $id
     * @return Response
     */
    public function destroy($id)
    {
        if ($this->service->deleteField($id)) {
            return response()->json(['status' => 200, 'message' => 'The specified resource has been deleted.'],
                200);
        }

        return response()->json(['status' => 422, 'campaign' => ['The specified resource could not be deleted.']],
            422);
    }
}
