<?php

namespace newsletters\Http\Controllers\Api;

use Illuminate\Http\Request;
use newsletters\Http\Controllers\Controller;
use newsletters\Http\Requests;
use newsletters\Repositories\TemplateRepository;
use newsletters\Services\TemplateService;

class TemplateController extends Controller
{

    /**
     * @var TemplateService
     */
    private $service;

    public function __construct(TemplateService $service)
    {
        $this->middleware('auth.basic');

        $this->service = $service;
    }

    /**
     * Display a listing of the resource.
     *
     * @param Request $request
     * @param TemplateRepository $repository
     * @return Response
     */
    public function index(Request $request, TemplateRepository $repository)
    {
        $templates = $this->service->findAllTemplates($request->has('paginate'), 10, $repository);

        return response()->json($templates, 200);
    }

    /**
     * Store a newly created resource in storage.
     *
     * @param  Request $request
     * @return Response
     */
    public function store(Request $request)
    {
        //
    }

    /**
     * Display the specified resource.
     *
     * @param  int $id
     * @param TemplateRepository $repository
     * @return Response
     */
    public function show($id, TemplateRepository $repository)
    {
        $template = $this->service->findTemplate($id, $repository);
        if (isset($template)) {
            return response()->json($template, 200);
        }

        return response()->json(['status' => 404, 'message' => 'The specified resource does not exist.'], 404);
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
        //
    }

    /**
     * Remove the specified resource from storage.
     *
     * @param  int $id
     * @param TemplateService $service
     * @param TemplateRepository $repository
     * @return Response
     */
    public function destroy($id, TemplateService $service, TemplateRepository $repository)
    {
        if ($service->deleteUnusedTemplate($id, $repository)) {
            return response()->json(['status' => 200, 'message' => 'The specified resource has been deleted.'], 200);
        }

        return response()->json(['status' => 422, 'message' => ['The specified resource could not be deleted.']],
            422);
    }
}
