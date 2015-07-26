<?php

namespace newsletters\Http\Controllers\Api;

use Illuminate\Database\Eloquent\ModelNotFoundException;
use Illuminate\Http\Request;

use newsletters\Http\Requests;
use newsletters\Http\Controllers\Controller;
use newsletters\Repositories\TemplateRepository;

class TemplateController extends Controller
{

    /**
     * @var TemplateRepository
     */
    private $repository;

    public function __construct(TemplateRepository $repository)
    {
        $this->middleware('auth');

        $this->repository = $repository;
    }

    /**
     * Display a listing of the resource.
     *
     * @param Request $request
     * @return Response
     */
    public function index(Request $request)
    {
        if($request->has('paginate')) {
            $perPage = ($request->has('per_page')) ? $request->input('per_page') : 15;
            $templates = $this->repository->paginate($perPage);
        } else {
            $templates = $this->repository->all();
        }
        return response()->json($templates, 200);
    }

    /**
     * Store a newly created resource in storage.
     *
     * @param  Request  $request
     * @return Response
     */
    public function store(Request $request)
    {
        //
    }

    /**
     * Display the specified resource.
     *
     * @param  int  $id
     * @return Response
     */
    public function show($id)
    {
        try {
            $template = $this->repository->find($id);

            return response()->json($template, 200);
        } catch (ModelNotFoundException $e) {
            return response()->json(['status' => 404, 'message' => 'The specified resource does not exist.'], 404);
        }
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
        //
    }

    /**
     * Remove the specified resource from storage.
     *
     * @param  int  $id
     * @return Response
     */
    public function destroy($id)
    {
        //
    }
}
