<?php

namespace newsletters\Http\Controllers\Api;

use Illuminate\Http\Request;
use newsletters\Http\Controllers\Controller;
use newsletters\Http\Requests\StoreTemplateRequest;
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
     * @return Response
     */
    public function index(Request $request)
    {
        $perPage = ($request->has('per_page')) ? $request->input('per_page') : 10;

        if ($request->has('paginate')) {
            $templates = $this->service->findAllTemplatesPaginated($perPage);
        } else {
            $templates = $this->service->findAllTemplates();
        }

        return response()->json($templates, 200);
    }

    /**
     * Store a newly created resource in storage.
     *
     * @param StoreTemplateRequest $request
     * @return Response
     */
    public function store(StoreTemplateRequest $request)
    {
        $template = $this->service->createTemplate($request->all());

        if (isset($template)) {
            return response()->json(['template' => $template->id], 200);
        }

        return response()->json(['message' => ['The specified resource could not be created.']], 412);
    }

    /**
     * Display the specified resource.
     *
     * @param  int $id
     * @return Response
     */
    public function show($id)
    {
        $template = $this->service->findTemplate($id);

        if (isset($template)) {
            return response()->json($template, 200);
        }

        return response()->json(['message' => ['The specified resource does not exist.']], 404);
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
        $template = $this->service->updateTemplate($request->all(), $id);

        if (isset($template)) {
            return response()->json(['template' => $template->id], 200);
        }

        return response()->json(['message' => ['The specified resource could not be updated.']],
            412);
    }

    /**
     * Remove the specified resource from storage.
     *
     * @param  int $id
     * @return Response
     */
    public function destroy($id)
    {
        if ($this->service->deleteUnusedTemplate($id)) {
            return response()->json(['message' => ['The specified resource has been deleted.']], 200);
        }

        return response()->json(['message' => ['The specified resource could not be deleted.']],
            422);
    }

    /**
     * Display the content of the specified template
     *
     * @param $id
     * @return \Illuminate\Contracts\Routing\ResponseFactory|\Illuminate\Http\JsonResponse|\Symfony\Component\HttpFoundation\Response
     */
    public function showContent($id)
    {
        $template = $this->service->findTemplate($id);
        if (isset($template)) {
            return response($template->content, 200);
        }

        abort(404);
    }
}
